// /opt/axentx/social-assembly/frontend/src/pages/DiscoverPage.jsx
import React, { useState } from 'react';
import CommunityCard from '../components/CommunityCard';
import { communitiesData } from '../data/communities'; // Assume this is a placeholder data file

const DiscoverPage = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const communitiesPerPage = 10;

  const indexOfLastCommunity = currentPage * communitiesPerPage;
  const indexOfFirstCommunity = indexOfLastCommunity - communitiesPerPage;
  const currentCommunities = communitiesData.slice(indexOfFirstCommunity, indexOfLastCommunity);

  const filteredCommunities = currentCommunities.filter((community) => {
    return (
      community.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      community.tags.some((tag) => tag.toLowerCase().includes(searchTerm.toLowerCase())) ||
      community.memberCount.toString().includes(searchTerm)
    );
  });

  const paginate = (pageNumber) => setCurrentPage(pageNumber);

  return (
    <div>
      <input
        type="text"
        placeholder="Search by name, tags, or member count"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
      />
      <div className="community-list">
        {filteredCommunities.map((community) => (
          <CommunityCard key={community.id} {...community} />
        ))}
      </div>
      <div className="pagination">
        {[...Array(Math.ceil(communitiesData.length / communitiesPerPage))].map((_, i) => (
          <button key={i} onClick={() => paginate(i + 1)}>
            {i + 1}
          </button>
        ))}
      </div>
    </div>
  );
};

export default DiscoverPage;